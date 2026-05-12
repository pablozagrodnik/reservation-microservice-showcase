import { apiPost } from './client';

export interface CreateReservationPayload {
    screening_id: number;
    seat_id: number;
    user_email: string;
}

export interface Reservation {
    id: number;
    screening_id: number;
    seat_id: number;
    user_email: string;
}

export function createReservation(
    payload: CreateReservationPayload,
    signal?: AbortSignal
): Promise<Reservation> {
    return apiPost<Reservation>('/reservations', payload, { signal });
}
