
export interface Heartbeat {
    id: number,
    service_id: number,
    response_time: number,
    status_code: number,
    is_success: boolean,
    created_at: Date,
}
