<script lang="ts">
	import {
		Navbar,
		NavbarBrand,
		Button,
		Modal,
		ModalBody,
		ModalFooter,
		ModalHeader,
		FormGroup,
		Label,
		Input
	} from 'sveltestrap';
	import type { DescanClient } from '../lib/api';
	import { onMount } from 'svelte';
	import Search from './Search.svelte';

	export let client: DescanClient;
  export let base: string;
	let alive = false;
	let message = '';
	let endpoint: string;

	client.alive.subscribe((val) => {
		alive = val;
	});
	client.message.subscribe((val) => {
		message = val;
	});

	onMount(async () => {
		endpoint = client.getEndpoint();
	});

	let open = false;

	function toggle() {
		endpoint = client.getEndpoint();
		open = !open;
	}

	async function saveChanges() {
		window.localStorage.setItem('endpoint', endpoint);
		console.log(endpoint);
		client.setEndpoint(endpoint);
		toggle();
	}

	async function cancelChanges() {
		toggle();
	}
</script>

<Navbar color="light" light expand="md">
	<NavbarBrand href="{base}/">DScan</NavbarBrand>
	<Search on:search />
	<Button on:click={toggle}>
		{#if alive}
			Connected
		{:else}
			Disconnected
		{/if}
	</Button>
	<Modal isOpen={open} {toggle}>
		<ModalHeader {toggle}>Status</ModalHeader>
		<ModalBody>
			{#if alive}
				Connected to indexer server. <br />
				{message}
			{:else}
				Not connected to indexer server.
			{/if}
			<hr />
			<FormGroup>
				<Label>Indexer URI</Label>
				<Input bind:value={endpoint} />
			</FormGroup>
		</ModalBody>
		<ModalFooter>
			<Button color="primary" on:click={saveChanges}>Save</Button>
			<Button color="secondary" on:click={cancelChanges}>Cancel</Button>
		</ModalFooter>
	</Modal>
</Navbar>
