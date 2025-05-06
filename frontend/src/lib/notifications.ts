import { writable, derived, type Writable, type Readable } from "svelte/store";

type NotificationType = {
	id: string;
	type: string;
	message: string;
	timeout: number;
};

function createNotificationStore() {
	const _notifications: Writable<Array<NotificationType>> = writable([]);

	function send(message: string, type: string = "default", timeout: number) {
		_notifications.update((state) => {
			return [...state, { id: id(), type, message, timeout }];
		});
	}

	const notifications: Readable<Array<NotificationType>> = derived(
		_notifications,
		($_notifications, set) => {
			set($_notifications);
			if ($_notifications.length > 0) {
				const timer = setTimeout(() => {
					_notifications.update((state) => {
						state.shift();
						return state;
					});
				}, $_notifications[0].timeout);
				return () => {
					clearTimeout(timer);
				};
			}
		}
	);
	const { subscribe } = notifications;

	return {
		subscribe,
		send,
		default: (msg: string, timeout: number) =>
			send(msg, "default", timeout),
		error: (msg: string, timeout: number) => send(msg, "error", timeout),
		info: (msg: string, timeout: number) => send(msg, "info", timeout),
		success: (msg: string, timeout: number) =>
			send(msg, "success", timeout),
	};
}

function id() {
	return "_" + Math.random().toString(36).substring(2, 9);
}

export const notifications = createNotificationStore();

//@ts-expect-error
window.notifications = notifications
