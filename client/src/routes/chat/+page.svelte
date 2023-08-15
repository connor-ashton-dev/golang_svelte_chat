<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { name, room } from '$lib/store';
	import { v4 as uuidv4 } from 'uuid';
	let messageFeed: string[] = [''];
	let currentUsers: string[] = [''];
	let message: string = '';
	let channelInputVal: string = '';
	let currentChannel = $room;
	let myId: string = $name + '::' + uuidv4();

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
		// ws = new WebSocket('ws://localhost:3000/ws/' + currentChannel);
		ws.onmessage = (event) => {
			const data = event.data;
			if (data.includes('users::')) {
				const userArray = data.split('::');
				console.log(userArray);
				currentUsers = userArray.slice(1, userArray.length);
			}
			if (data.includes('message::')) {
				const user = data.split('::')[1];
				const message = data.split('::')[2];
				messageFeed = [...messageFeed, user + ': ' + message];
			}
		};
		ws.onopen = () => {
			if ($name != '') {
				ws.send('join::' + myId);
			}
		};
	};

	const changeWS = () => {
		ws.close();
		currentChannel = channelInputVal;
		messageFeed = [];
		setupWS();
	};

	const sendMessage = () => {
		ws.send('message::' + $name + '::' + message);
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
		<p>Current Users:</p>
		<ul class="text-green-500">
			{#each currentUsers as user}
				<li>{user}</li>
			{/each}
		</ul>
	</div>

	<div class=" flex flex-col justify-center items-center mt-20">
		<p class="text-4xl mb-10">Convo:</p>
		{#each messageFeed as message}
			<p>{message}</p>
		{/each}
		<div class="w-full flex flex-row justify-center mt-10 lg:mb-60 sm:mb-10">
			<input
				class="border border-black outline-none p-1 rounded-l-lg"
				placeholder="Send Message"
				bind:value={message}
			/>
			<button class="text-white bg-blue-500 p-2 rounded-r-lg" on:click={sendMessage}>Send</button>
		</div>
	</div>
</div>
