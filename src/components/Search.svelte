<script lang="ts">
	let query: string;
	function handleSearch() {
		switch (true) {
			case /0x[0-9a-fA-F]{64}/.test(query):
				window.location.href = `/transaction?id=${query}`;
				break;
			case /0x[0-9a-fA-F]{40}/.test(query):
				window.location.href = `/address?id=${query}`;
				break;
			case /[0-9]*/.test(query):
				window.location.href = `/block?id=${query}`;
				break;
			default:
				break;
		}
	}
</script>

<form on:submit|preventDefault={handleSearch}>
	<div class="input-group input-group-lg">
		<input
			type="text"
			class="form-control"
			placeholder="Address, block number, or transaction hash"
			aria-label="Address, block number, or transaction hash"
			aria-describedby="basic-addon2"
			bind:value={query}
		/>
		<div class="input-group-append">
			<button class="btn btn-lg btn-primary" type="button" on:click={handleSearch}>Search</button>
		</div>
	</div>
</form>
