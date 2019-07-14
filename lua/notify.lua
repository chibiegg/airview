local NOTIFY_URL = "http://192.168.0.11:8000/notify/"   -- サーバのアドレスに設定する

local PHOTO_TARGET_PATH = "/DCIM"

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


function find_latest_directory()
  local latest_dir = nil
  local latest_number = 0
  for dir in lfs.dir(PHOTO_TARGET_PATH) do
    c = string.match(dir, "^(%d%d%d)")
    print(c)
    if c ~= nil then
      local n = tonumber(c)
      if n > latest_number then
        latest_number = n
        latest_dir = dir
      end
    end
  end
  return PHOTO_TARGET_PATH.."/"..latest_dir
end



local photo_dir = find_latest_directory()
local filename = find_latest_path(photo_dir)
local current_number = get_file_number(filename) - 1

print(photo_dir)

while true do
  for file in lfs.dir(photo_dir) do
      local path = photo_dir.."/"..file
      local suffix = string.sub(string.lower(path), -3)
      local number = get_file_number(path)
      if suffix == "jpg" and number > current_number then
        print(path)
        post_text(NOTIFY_URL, path)
        current_number = number
      end
  end

  sleep(500)
end
