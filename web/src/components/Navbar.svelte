<script lang="ts">
    import "../app.css";

	import { Navigation } from '@skeletonlabs/skeleton-svelte';
    import { ShoppingCart, User } from '@lucide/svelte';
</script>

<header class="border-b border-surface-200-800 bg-surface-100-900">
    <div class="mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 items-center justify-between">
            <div class="mx-12">
                <a href="/" class="flex items-center gap-2 font-semibold">
                    <img src="/assets/logo.svg" alt="Stocklet" class="h-8 w-auto dark:hidden" />
                    <img src="/assets/logo-light.svg" alt="Stocklet" class="h-8 w-auto hidden dark:block" />
                </a>
            </div>

            <Navigation.Bar classes="bg-transparent w-auto p-0 mx-12">
                <Navigation.Tile href="/" label="Home" />
                <Navigation.Tile href="/about" label="About" />

                <div class="flex flex-row gap-4">
                    <Navigation.Tile href="#" labelExpanded="Cart" title="Cart">
                        <ShoppingCart class="h-6 w-6" />
                    </Navigation.Tile>
                    <Navigation.Tile href="#" labelExpanded="User" title="User">
                        <User class="h-6 w-6" />
                    </Navigation.Tile>
                </div>
            </Navigation.Bar>
        </div>
    </div>
</header>