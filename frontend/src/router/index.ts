import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import RoomView from '@/views/RoomView.vue';
import CheckoutView from '@/views/CheckoutView.vue';

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
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

export default router;