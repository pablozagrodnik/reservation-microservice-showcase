export interface Room {
    id: number;
    name: string;
}

export interface Screening {
    id: number;
    start_time: string;
    room_id: number;
    room: Room;
}

export interface Movie {
    id: number;
    title: string;
    description: string;
    poster: string;
    screenings: Screening[];
}

export interface Seat {
    id: number;
    row: number;
    col: number;
    is_taken: boolean;
}