<script>
	import { goto } from '$app/navigation';
	import { name, room } from '$lib/store';

	const handleClick = () => {
		if ($name != '') {
			if ($room == '') {
				room.set('global');
			}
			goto('/chat');
		} else {
			alert('Please enter your name');
		}
	};
</script>

<div class="w-screen h-screen flex flex-col items-center justify-center">
	<h1 class="text-3xl pb-4">Welcome {$name}</h1>
	<input
		class="border border-black outline-none p-1 rounded-lg"
		type="text"
		placeholder="name here..."
		bind:value={$name}
	/>
	<input
		class="border border-black outline-none p-1 w-60 rounded-lg mt-2"
		type="text"
		placeholder="room here... (empty for global)"
		bind:value={$room}
	/>

	<div class="flex flex-row gap-4 py-2">
		<button class="text-white bg-blue-500 p-2 rounded-lg" on:click={handleClick}
			>Start chatting</button
		>
		<a class="text-white bg-blue-500 p-2 rounded-lg" href="/feed">Live feed</a>
	</div>
</div>
