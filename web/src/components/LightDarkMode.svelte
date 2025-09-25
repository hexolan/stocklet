<script lang="ts">
  import { Label, Switch } from "bits-ui";

  let checked = $state(false);

  $effect(() => {
    const mode = localStorage.getItem('mode') || 'light';
    checked = mode === 'dark';
    document.documentElement.setAttribute('data-mode', mode);
  });

  const onCheckedChange = (event: { checked: boolean }) => {
    const mode = event.checked ? 'dark' : 'light';
    document.documentElement.setAttribute('data-mode', mode);
    localStorage.setItem('mode', mode);
    checked = event.checked;
  };
</script>

<svelte:head>
  <script>
    const mode = localStorage.getItem('mode') || 'light';
    document.documentElement.setAttribute('data-mode', mode);
  </script>
</svelte:head>

<div class="flex items-center space-x-3">
  <Switch.Root
    id="light-dark"
    checked={checked}
    onCheckedChange={(change) => { onCheckedChange({ checked: change }) }}
    class="focus-visible:ring-foreground focus-visible:ring-offset-background data-[state=checked]:bg-foreground data-[state=unchecked]:bg-dark-10 data-[state=unchecked]:shadow-mini-inset dark:data-[state=checked]:bg-foreground focus-visible:outline-hidden peer inline-flex h-[36px] min-h-[36px] w-[60px] shrink-0 cursor-pointer items-center rounded-full px-[3px] transition-colors focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
  >
    <Switch.Thumb
      class="bg-background data-[state=unchecked]:shadow-mini dark:border-background/30 dark:bg-foreground dark:shadow-popover pointer-events-none block size-[30px] shrink-0 rounded-full transition-transform data-[state=checked]:translate-x-6 data-[state=unchecked]:translate-x-0 dark:border dark:data-[state=unchecked]:border"
    />
  </Switch.Root>
  <Label.Root for="light-dark" class="text-sm font-medium">Light / Dark</Label.Root>
</div>