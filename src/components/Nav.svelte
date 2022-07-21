<script lang="ts">
	import { Navbar, NavbarBrand } from 'sveltestrap';

	import Search from './Search.svelte';
</script>

<Navbar color="light" light expand="md">
	<NavbarBrand href="/">DScan</NavbarBrand>
	<Search on:search />
</Navbar>
