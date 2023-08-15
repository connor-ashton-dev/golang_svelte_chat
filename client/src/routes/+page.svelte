<script lang="ts">
	import { onMount } from 'svelte';
	let messageFeed: string[] = [''];
	let message: string = '';
	let channelInputVal: string = 'global';
	let currentChannel: string = 'global';

	let name: string = '';
	let ws: WebSocket;
	onMount(() => {
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
		if (!name) {
			alert('Please enter your name');
			return;
		}
		ws.send(name + ': ' + message);
		message = '';
	};
</script>

<h1>Websocket chat</h1>

<p>Listening on channel: {currentChannel}</p>

<label for="channel">Channel:</label>
<input id="channel" placeholder="write channel here" bind:value={channelInputVal} />
<button on:click={changeWS}>Switch</button>

<label for="name">Name:</label>
<input id="name" placeholder="write name here" bind:value={name} />

<label for="message">Message:</label>
<input id="message" placeholder="write message here" bind:value={message} />
<button on:click={sendMessage}>Send</button>

{#each messageFeed as message}
	<p>{message}</p>
{/each}
