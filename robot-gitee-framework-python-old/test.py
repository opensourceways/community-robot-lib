# from dataclasses import dataclass
# import dacite
#
#
# @dataclass
# class Body(object):
#     day: int
#     month: int
#     year: int
#
#
# @dataclass
# class Dat(object):
#     response_time: int
#     body: Body
#
#
# data = {'response_time': 12, 'body': {'day': 1, 'month': 2, 'year': 3}}
#
# dat = dacite.from_dict(Dat, data)
# print(dat)
# print(dat.body.day)

class Solution:
    def xxx(self):
        print('可以把一些前置方法写在类里面或者类外面，供unlock调用')

    def unlock(self, init_state, dst_state):
        # 具体实现代码
        return 0


if __name__ == '__main__':
    init_state = '0023'
    dst_state = '0059'
    function = Solution()
    result = function.unlock(init_state, dst_state)
    print(result)
