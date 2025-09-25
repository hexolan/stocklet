<script lang="ts">
	import Navbar from '../components/Navbar.svelte';
	import Footer from '../components/Footer.svelte';

	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<title>Stocklet</title>

	<link rel="icon" href={favicon} />
</svelte:head>

<main data-theme="catppuccin" class="min-h-screen h-full">
	<div>
		<Navbar />

		{@render children?.()}

		<Footer />
	</div>
</main>
