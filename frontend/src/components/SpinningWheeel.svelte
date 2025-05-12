<div aria-busy="true">
	<span></span>
</div>

<style>
div {
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
}
span {
	--size: 8rem;
	height: var(--size);
	width: var(--size);
	border-color: var(--accent) transparent var(--accent) var(--accent);
	border-width: calc(var(--size) / 15);
	border-style: solid;
	border-image: initial;
	border-radius: 50%;
	animation: 1s linear 0s infinite normal none running rotate;
}
@keyframes rotate {
	0% {
		transform: rotate(0);
	}
	100% {
		transform: rotate(360deg);
	}
}
</style>
