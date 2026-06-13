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

export function createReservations(
    payload: CreateReservationPayload[],
    signal?: AbortSignal
): Promise<{ message: string }> {
    return apiPost<{ message: string }>('/reservations', payload, { signal });
}
