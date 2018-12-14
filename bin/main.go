package main

import (
        "fmt"
        "os"
        "net/http"
        "io"
        "io/ioutil"
        "path"
        "strings"
        "encoding/json"
        "sort"

        // "golang.org/x/net/websocket"

        "goji.io"
        "goji.io/pat"

        log "github.com/sirupsen/logrus"
        "github.com/golang-collections/go-datastructures/queue"
)

var q *queue.Queue
var logger *log.Entry
var latestFile string

func notify(w http.ResponseWriter, r *http.Request) {
        b, _ := ioutil.ReadAll(r.Body)
        hp := strings.Split(r.RemoteAddr, ":")
        u := fmt.Sprintf("http://%s%s", hp[0], string(b))
        q.Put(u)
        logger.Infof("Enqueue %s", u)
}

func download(u string, dir string) (string, error) {
  basename := path.Base(u)
  filename := fmt.Sprintf("%s/%s", dir, basename)

  response, err := http.Get(u)
  if err != nil {
    return "", err
  }
  defer response.Body.Close()

  file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
  if err != nil {
    return "", err
  }
  defer file.Close()

  _, err = io.Copy(file, response.Body)
  if err != nil {
    return "", err
  }

  return basename, nil
}

func download_loop(q *queue.Queue) {
  for {
    dl, _ := q.Get(1)
    for _, di := range dl{
        u := di.(string)
        logger.Debugf("Get %s", u)
        dist, err := download(u, "downloads")
        if err != nil {
          logger.Errorf("Download Error: %s", err)
        }else{
          logger.Infof("Download complete: %s", u)
          latestFile = dist
        }
    }
  }
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
  // TODO: WebSocketでの通知を実装する
}

func main() {
  // Log as JSON instead of the default ASCII formatter.
  //log.SetFormatter(&log.JSONFormatter{})

  // Output to stdout instead of the default stderr
  // Can be any io.Writer, see below for File example
  log.SetOutput(os.Stdout)

  // Only log the warning severity or above.
  log.SetLevel(log.DebugLevel)

  logger = log.WithFields(log.Fields{})

  q = queue.New(100)

  go download_loop(q)

  latestFile = "Initial.JPG"

  mux := goji.NewMux()
  mux.HandleFunc(pat.Get("/ws"), func(w http.ResponseWriter, r *http.Request) {
    serveWs(w, r)
  })
  mux.HandleFunc(pat.Get("/latest"), func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("/downloads/" + latestFile))
  })
  mux.HandleFunc(pat.Get("/list"), func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    files, err := ioutil.ReadDir("./downloads")
    if err != nil {
        logger.Error(err)
    }

    paths := make([]string, 0, len(files))

    for _, file := range files {
        if file.IsDir() {
            continue
        }
        if strings.HasSuffix(strings.ToLower(file.Name()), "jpg") == false {
          continue
        }
        paths = append(paths, "/downloads/" + file.Name())
    }

    sort.Strings(paths)
    b, _ := json.Marshal(paths)
    w.Write([]byte(b))
  })
  mux.HandleFunc(pat.Post("/notify/"), notify)
  mux.Handle(pat.Get("/downloads/*"), http.StripPrefix("/downloads", http.FileServer(http.Dir("downloads"))))
  mux.Handle(pat.Get("/*"), http.StripPrefix("/", http.FileServer(http.Dir("static"))))
  http.ListenAndServe("0.0.0.0:8000", mux)
}
