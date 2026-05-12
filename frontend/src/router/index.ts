import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import RoomView from '@/views/RoomView.vue';
import CheckoutView from '@/views/CheckoutView.vue';
import AdminLayout from '@/views/admin/AdminLayout.vue';
import AdminMovies from '@/views/admin/AdminMovies.vue';
import AdminRooms from '@/views/admin/AdminRooms.vue';
import AdminScreenings from '@/views/admin/AdminScreenings.vue';

const routes = [
    {
        path: '/',
        name: 'home',
        component: HomeView
    },
    {
        path: '/screening/:id',
        name: 'room',
        component: RoomView,
        props: true
    },
    {
        path: '/checkout',
        name: 'checkout',
        component: CheckoutView
    },
    {
        path: '/admin',
        component: AdminLayout,
        children: [
            { path: '', redirect: { name: 'admin-movies' } },
            { path: 'movies', name: 'admin-movies', component: AdminMovies },
            { path: 'rooms', name: 'admin-rooms', component: AdminRooms },
            { path: 'screenings', name: 'admin-screenings', component: AdminScreenings }
        ]
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

export default router;
