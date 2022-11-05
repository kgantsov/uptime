import {
    Color,
   } from "@tremor/react";

export interface Heartbeat {
    id: number,
    service_id: number,
    response_time: number,
    status_code: number,
    status: string,
    created_at: Date,
}

interface StatusColors {
    [name: string]: Color
}

export const STATUS_COLORS_MAP: StatusColors = {
    "UP": "emerald",
    "DOWN": "rose",
    "TIMEOUT": "amber",
    "FAILED": "rose",
    "UNKNOWN": "gray",
}
