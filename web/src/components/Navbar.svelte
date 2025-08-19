<script lang="ts">
	import { Navigation } from '@skeletonlabs/skeleton-svelte';
</script>

<header class="border-b border-surface-200-800 bg-surface-100-900">
    <div class="mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 items-center justify-between">
            <a href="/" class="flex items-center gap-2 font-semibold">
                <img src="/assets/logo.svg" alt="Stocklet" class="h-8 w-auto dark:hidden" />
                <img src="/assets/logo-light.svg" alt="Stocklet" class="h-8 w-auto hidden dark:block" />
            </a>
            <Navigation.Bar classes="bg-transparent width-auto h-12 p-0">
                <Navigation.Tile href="/" label="Home" />
                <Navigation.Tile href="#" label="Products" />
                <Navigation.Tile href="#" label="Orders" />
                <Navigation.Tile href="#" label="Account" />
            </Navigation.Bar>
        </div>
    </div>
</header>