<script lang="ts">
	import { onMount } from 'svelte';
	let messageFeed: string[] = [''];

	let ws: WebSocket;
	onMount(() => {
		setupWS();
	});

	const setupWS = () => {
		ws = new WebSocket('wss://server-xcjt64gdjq-wm.a.run.app/ws/feed');
		ws.onmessage = (event) => {
			messageFeed = [...messageFeed, event.data];
		};
	};
</script>

<div class="p-6">
	<h1>Live Feed</h1>

	{#each messageFeed as message}
		<p>{message}</p>
	{/each}
</div>
