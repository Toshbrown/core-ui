<template>
	<iframe :src="url"></iframe>
</template>
<script>
	export default {
		name: 'appUI',
		props: ['app', 'path'],
		computed: {
			url() {
				let url = '/' + this.app + '/ui';
				if (this.path) {
					url = url + '/' + this.path;
				}

				let search = window.location.search;
				if (search === '') {
					if (this.isMobile) {
						url = url + '?mobile=true';
					}
				} else {
					url = url + window.location.search;
					if (this.isMobile) {
						url = url + '&mobile=true';
					}
				}
				return url;
			}
		},
		mounted() {
			this.$parent.title = this.app;
			this.$parent.backRoute = "/";
		},
	}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
	iframe {
		width: 100%;
		height: 100%;
		border: none;
	}
</style>