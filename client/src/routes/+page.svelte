<script lang="ts">
	import { onMount } from 'svelte';
	let messageFeed: string[] = [''];
	let myMessage: string = '';
	let selectedChannel: string = 'global';
	let currentChannel: string = 'global';

	let name: string = '';
	let ws: WebSocket;
	onMount(() => {
		ws = new WebSocket('ws://localhost:3000/ws/' + selectedChannel);
		ws.onmessage = (event) => {
			messageFeed = [...messageFeed, event.data];
		};
	});

	const changeWS = () => {
		ws.close();
		ws = new WebSocket('ws://localhost:3000/ws/' + selectedChannel);
		ws.onmessage = (event) => {
			messageFeed = [...messageFeed, event.data];
		};
		currentChannel = selectedChannel;
		messageFeed = [];
	};

	const sendMessage = () => {
		if (!name) {
			alert('Please enter your name');
			return;
		}
		ws.send(name + ': ' + myMessage);
		myMessage = '';
	};
</script>

<h1>Websocket chat</h1>

<p>Listening on channel: {currentChannel}</p>

<label for="channel">Channel:</label>
<input id="channel" placeholder="write channel here" bind:value={selectedChannel} />
<button on:click={changeWS}>Switch</button>

<label for="name">Name:</label>
<input id="name" placeholder="write name here" bind:value={name} />

<label for="message">Message:</label>
<input id="message" placeholder="write message here" bind:value={myMessage} />
<button on:click={sendMessage}>Send</button>

{#each messageFeed as message}
	<p>{message}</p>
{/each}
