
export interface Notification {
  id: number,
  name: string,
  callback: string,
  callback_chat_id: string,
  callback_type: string,
  created_at: string,
  updated_at: string,
}

export interface Service {
  id: number,
  name: string,
  url: string,
  enabled: boolean,
}
