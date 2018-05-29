local PHOTO_TARGET_PATH = "/DCIM/100__TSB"
local NOTIFY_URL = "http://192.168.0.11:8000/notify/"

function find_latest_path(basedir)
    local max_mod = 0
    local max_mod_path = nil
    for file in lfs.dir(basedir) do
        local path = basedir.."/"..file
        local last_mod = lfs.attributes(path, "modification")
        if last_mod > max_mod then
            max_mod = last_mod
            max_mod_path = path
        end
    end
    return max_mod_path
end

function post_text(url, payload)
  fa.request{
    url = url,
    method = "POST",
    headers = {
        ["Content-Length"] = string.len(payload),
        ["Content-Type"] = "text/plain"
    },
    body = payload
  }
end

function get_file_number(filename)
  return tonumber(string.sub(string.lower(filename), -8, -5), 10)
end

filename = find_latest_path(PHOTO_TARGET_PATH)
local current_number = get_file_number(filename) - 1

while true

  for file in lfs.dir(PHOTO_TARGET_PATH) do
      local path = PHOTO_TARGET_PATH.."/"..file
      local suffix = string.sub(string.lower(path), -3)
      local number = get_file_number(path)
      if suffix == "jpg" and number > current_number then
        post_text(NOTIFY_URL, path)
        current_number = number
      end
  end

  sleep(1000)
do
