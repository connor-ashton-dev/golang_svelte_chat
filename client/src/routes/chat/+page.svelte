<script lang="ts">
	import { onMount } from 'svelte';
	import { name } from '$lib/store';
	let messageFeed: string[] = [''];
	let message: string = '';
	let channelInputVal: string = '';
	let currentChannel: string = 'global';

	let ws: WebSocket;
	onMount(() => {
		if ($name == '') {
			// redirect back to homepage
			window.location.href = '/';
		}
		setupWS();
	});

	const setupWS = () => {
		ws = new WebSocket('wss://server-xcjt64gdjq-wm.a.run.app/ws/' + currentChannel);
		ws.onmessage = (event) => {
			messageFeed = [...messageFeed, event.data];
		};
	};

	const changeWS = () => {
		ws.close();
		currentChannel = channelInputVal;
		messageFeed = [];
		setupWS();
	};

	const sendMessage = () => {
		ws.send($name + ': ' + message);
		message = '';
	};
</script>

<div class="w-screen h-screen p-4">
	<div class="w-screen flex flex-col items-center justify-center">
		<p class="text-3xl py-4">Listening on channel: {currentChannel}</p>
		<div class="w-full flex flex-row justify-center">
			<input
				class="border border-black outline-none p-1 rounded-l-lg"
				id="channel"
				placeholder="Change channel"
				bind:value={channelInputVal}
			/>
			<button class="text-white bg-blue-500 p-2 rounded-r-lg" on:click={changeWS}>Change</button>
		</div>
	</div>

	<div class="w-full flex flex-col justify-center items-center mt-20">
		<p class="text-4xl">Convo:</p>
		<div class="py-4">
			{#each messageFeed as message}
				<p>{message}</p>
			{/each}
		</div>
		<div class="w-full flex flex-row justify-center mb-60">
			<input
				class="border border-black outline-none p-1 rounded-l-lg"
				placeholder="Send Message"
				bind:value={message}
			/>
			<button class="text-white bg-blue-500 p-2 rounded-r-lg" on:click={sendMessage}>Change</button>
		</div>
	</div>
</div>
