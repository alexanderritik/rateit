const KEY = 'catch_it_device_id'

export function getOrCreateDeviceID(): string {
  let id = localStorage.getItem(KEY)
  if (!id) {
    id = crypto.randomUUID()
    localStorage.setItem(KEY, id)
  }
  return id
}
