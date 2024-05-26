

export function SIToBytes(memory: string): number {
	const regex = /^([0-9]+)(k|m|g|ki|mi|gi)$/i;
	const matched = regex.test(memory.toLowerCase());
	if (!matched) {
		return -1;
	}

	const amount = parseInt(memory.match(/^[0-9]+/)![0], 10)
	if (isNaN(amount)) {
		return -1;
	}

	const unit = memory.match(/(k|m|g|ki|mi|gi)$/i)![0]; // Type assertion for capture group
	let multiplier: number;
	switch (unit) {
		case "k":
		case "ki":
			multiplier = 1024;
			break;
		case "m":
		case "mi":
			multiplier = 1024 * 1024;
			break;
		case "g":
		case "gi":
			multiplier = 1024 * 1024 * 1024;
			break;
		default:
			return -1
	}

	return amount * multiplier;
}

export function byteToGi(memory: number) {
	return memory / 1024 / 1024 / 1024
}