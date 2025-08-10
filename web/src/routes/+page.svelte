<script lang="ts">
	import { Navigation, Avatar, Rating } from '@skeletonlabs/skeleton-svelte';
	import Footer from '../components/Footer.svelte';
</script>

<div data-theme="theme" class="min-h-screen flex flex-col">
	<!-- Navbar -->
	<header class="border-b border-surface-200-800 bg-surface-100-900">
		<div class="mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-16 items-center justify-between">
				<a href="/" class="flex items-center gap-2 font-semibold">
					<img src="/assets/vector/default-monochrome.svg" alt="Stocklet" class="h-8 w-auto" />
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

	<!-- Hero -->
	<section class="relative overflow-hidden bg-surface-50-950">
		<div class="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8">
			<div class="grid items-center gap-12 md:grid-cols-2">
				<div>
					<h1 class="text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl">
						Phasellus ut pulvinar nisl.
					</h1>
					<p class="mt-6 max-w-prose text-lg opacity-80">
						Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer dignissim odio orci, in volutpat nunc lacinia et. Phasellus ut dictum nibh. Donec hendrerit nunc sed elit porttitor, in efficitur arcu hendrerit.
                        Aenean erat est, maximus id congue ut, venenatis eu nibh. Vivamus vitae varius erat, eu tempor urna. Aenean viverra sed erat sed viverra. Fusce vulputate, mauris eget hendrerit venenatis, eros felis rutrum orci, in commodo mauris elit sit amet arcu.
					</p>
					<div class="mt-8 flex flex-wrap items-center gap-3">
						<a href="#get-started" class="rounded-container preset-filled-primary-500 px-5 py-3 text-sm">Start browsing</a>
						<a href="/docs" class="rounded-container border border-surface-300-700 px-5 py-3 text-sm hover:bg-surface-100-900">Check orders</a>
					</div>
				</div>
				<div class="relative">
					<div class="aspect-[4/3] w-full overflow-hidden rounded-container border border-surface-200-800 bg-surface-50-950 shadow-sm">
						<div class="h-full w-full opacity-70" style="background: conic-gradient(at 20% 20%, var(--color-primary-50), transparent, var(--color-primary-50));"></div>
					</div>
					<div class="pointer-events-none absolute -right-10 -top-10 h-40 w-40 rounded-full preset-filled-secondary-500 blur-2xl opacity-20"></div>
				</div>
			</div>
		</div>
	</section>

	<!-- Testimonials -->
	<section id="testimonials" class="bg-surface-50-950">
		<div class="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
			<h2 class="text-center text-3xl font-bold tracking-tight sm:text-4xl">Loved by providers</h2>
			<p class="mx-auto mt-3 max-w-2xl text-center opacity-80">
				Providers use Stocklet to sell their products with confidence.
			</p>
			<div class="mt-10 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
				<figure class="rounded-container border border-surface-200-800 bg-surface-50-950 p-6 shadow-sm">
					<blockquote>“Vivamus sit amet est eu urna lobortis cursus sed at odio.”</blockquote>
					<figcaption class="mt-4 flex items-center gap-3 text-sm opacity-70">
						<Avatar name="John Doe" classes="h-8 w-8" />
						<div>
							<div class="font-medium">John Doe</div>
							<div>CEO, Acme Inc.</div>
						</div>
						<Rating value={5} readOnly controlGap="gap-1" />
					</figcaption>
				</figure>
				<figure class="rounded-container border border-surface-200-800 bg-surface-50-950 p-6 shadow-sm">
					<blockquote>“Aliquam velit augue, varius quis est ac, interdum tincidunt lorem.”</blockquote>
					<figcaption class="mt-4 flex items-center gap-3 text-sm opacity-70">
						<Avatar name="Jane Doe" classes="h-8 w-8" />
						<div>
							<div class="font-medium">Jane Doe</div>
							<div>Product Lead, Wayne Industries.</div>
						</div>
						<Rating value={5} readOnly controlGap="gap-1" />
					</figcaption>
				</figure>
				<figure class="rounded-container border border-surface-200-800 bg-surface-50-950 p-6 shadow-sm">
					<blockquote>“Aliquam a lobortis nisi. Nullam varius a diam nec vulputate.”</blockquote>
					<figcaption class="mt-4 flex items-center gap-3 text-sm opacity-70">
						<Avatar name="Alice Wonderland" classes="h-8 w-8" />
						<div>
							<div class="font-medium">Alice Wonderland</div>
							<div>Managing Directory, E Corp.</div>
						</div>
						<Rating value={5} readOnly controlGap="gap-1" />
					</figcaption>
				</figure>
			</div>
		</div>
	</section>

	<Footer />
</div>
