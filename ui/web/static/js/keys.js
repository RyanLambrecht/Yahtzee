document.addEventListener('keydown', (e) => {
	if (e.repeat || e.ctrlKey || e.metaKey || e.altKey) return;
	if (e.target.matches('input[type="text"], input:not([type]), textarea, select')) return;

	const el = document.querySelector(`[data-key="${e.key}"]`);
	if (!el) return;

	e.preventDefault();
	document.activeElement?.blur();
	el.click();
});
