// 这里的逻辑是用来实现一个加法器
var Calculator = function() {
    var numbers = [],
        resultSum = 0,
        isOverflow = false,
        dotDigit = 0; // 小数位；i点击逗号的时候变成1，之后再点击数字就累加
    var init = function() {
        _events();
    };
    var getSum = function() {
        return parseFloat(resultSum);
    };
    var _events = function() {
        // 计算按钮
        $('#keyZone').find('button').on('tap', function(e) {

            e.preventDefault();
            var $this = $(this),
                $numberZone = $('#numberZone'),
                $resultZone = $('#resultZone'),
                tag = $this.attr('data-tag'),
                l = numbers.length;
            switch (tag) {
                case 'clear':
                    numbers = [];
                    dotDigit = 0;
                    resultSum = 0;
                    isOverflow = false;
                    break;
                case 'delete':
                    if (l > 0) {
                        if (dotDigit > 0) {
                            dotDigit--;
                        }
                        numbers.pop();
                    }
                    break;
                case 'plus':
                    _pushNumber('+');
                    break;
                case 'scan':
                    // TODO 调用系统的摄像头
                    break;
                case 'dot':
                    _pushNumber('.');
                    break;
                default:
                    if (isOverflow) {
                        return;
                    }
                    _pushNumber(tag);
            }
            l = numbers.length;
            if (l === 0) {
                $numberZone.text('0');
                $resultZone.text('= 0');
                return;
            }
            var expression = numbers.join('');
            $numberZone.text(expression);
            resultSum = 0;

            var nums = expression.split('+');
            for (var i = 0, len = nums.length; i < len; i++) {
                if (nums[i] === '') {
                    continue;
                }
                var a = parseFloat(nums[i]);
                resultSum = accAdd(resultSum, a);
            }
            resultSum = resultSum.toFixed(2);
            $resultZone.text('= ' + resultSum);

            isOverflow = false;
            if (resultSum.length > 10) {
                isOverflow = true;
                window.alert("输入值已经最大了");
            }
        });
    };
    // 点击键盘上的数字的事件
    var _pushNumber = function(num) {
        var l = numbers.length;
        switch (num) {
            case '.':
                if (l === 0 || numbers[l - 1] === '+') {
                    numbers[l] = '0';
                    l++;
                }
                if (dotDigit > 0) {
                    return;
                }
                dotDigit = 1;
                break;
            case '+':
                if (l === 0) {
                    return;
                }
                if (numbers[l - 1] === '+') {
                    return;
                }
                if (numbers[l - 1] === '.') {
                    l--;
                }
                dotDigit = 0;
                break;
            default:
                if (dotDigit > 3) {
                    return;
                }
                if (dotDigit > 0) {
                    dotDigit++;
                }
        }
        numbers[l] = num;
    };
    // 精确加法
    var accAdd = function(arg1, arg2) {
        var r1, r2, m, c;
        try {
            r1 = arg1.toString().split(".")[1].length;
        } catch (e) {
            r1 = 0;
        }
        try {
            r2 = arg2.toString().split(".")[1].length;
        } catch (e) {
            r2 = 0;
        }
        c = Math.abs(r1 - r2);
        m = Math.pow(10, Math.max(r1, r2));
        if (c > 0) {
            var cm = Math.pow(10, c);
            if (r1 > r2) {
                arg1 = Number(arg1.toString().replace(".", ""));
                arg2 = Number(arg2.toString().replace(".", "")) * cm;
            } else {
                arg1 = Number(arg1.toString().replace(".", "")) * cm;
                arg2 = Number(arg2.toString().replace(".", ""));
            }
        } else {
            arg1 = Number(arg1.toString().replace(".", ""));
            arg2 = Number(arg2.toString().replace(".", ""));
        }
        return (arg1 + arg2) / m;
    };
    return {
        init: init,
        getSum: getSum
    };
}(window);
