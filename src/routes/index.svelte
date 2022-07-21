<script lang="ts">
	import { Col, Container, Row } from 'sveltestrap';
	import Search from '../components/Search.svelte';
</script>

<svelte:head>
	<title>Decentralized Explorer</title>
</svelte:head>

<Container>
	<div class="padded" />
	<Row>
		<Col><h1>Decentralized Explorer</h1></Col>
	</Row>
	<Row>
		<Col>
			<Search />
		</Col>
	</Row>
</Container>

<style>
	.padded {
		margin-top: 2rem;
	}
</style>
