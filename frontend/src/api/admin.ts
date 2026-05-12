import type { Room } from '@/types';
import { apiGet, apiPost } from './client';

export interface CreateMoviePayload {
    title: string;
    description: string;
    poster: string;
}

export interface AdminMovie {
    id: number;
    title: string;
    description: string;
    poster: string;
}

export function createMovie(
    payload: CreateMoviePayload,
    signal?: AbortSignal
): Promise<AdminMovie> {
    return apiPost<AdminMovie>('/admin/movies', payload, { signal });
}

export interface CreateRoomPayload {
    name: string;
    rows: number;
    cols: number;
}

export interface CreateRoomResponse {
    message: string;
    room: {
        id: number;
        name: string;
    };
}

export function createRoom(
    payload: CreateRoomPayload,
    signal?: AbortSignal
): Promise<CreateRoomResponse> {
    return apiPost<CreateRoomResponse>('/admin/rooms', payload, { signal });
}

export function fetchRooms(signal?: AbortSignal): Promise<Room[]> {
    return apiGet<Room[]>('/admin/rooms', { signal });
}

export interface CreateScreeningPayload {
    movie_id: number;
    room_id: number;
    start_time: string;
}

export interface AdminScreening {
    id: number;
    room_id: number;
    start_time: string;
}

export function createScreening(
    payload: CreateScreeningPayload,
    signal?: AbortSignal
): Promise<AdminScreening> {
    return apiPost<AdminScreening>('/admin/screenings', payload, { signal });
}
