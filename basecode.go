package tghost

const baseCode = `
function blankChecker(_) {
	return "";
}

function mustEqualChecker(value) {
	return input => (input == value ? '' : '请输入' + value);
}

function mustInChecker(values) {
	return input => (values.includes(input) ? '' : '请输入' + values.join('或'));
}

function numberInRangeChecker(min, max, minErrMsg, maxErrMsg) {
	return (input) => {
		let x = parseInt(input)
		if(!Number.isInteger(x)) {
			return "请输入整数！";
		}
		if(x < min) {
			if(!minErrMsg) {
				return "输入最小值为" + min + "！";
			}
			return minErrMsg;
		}
		if(x > max) {
			if(!maxErrMsg) {
				return "输入最大值为" + max + "！";
			}
			return maxErrMsg;
		}
		return "";
	}
}
`
