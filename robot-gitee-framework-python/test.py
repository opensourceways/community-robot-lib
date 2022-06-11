import json


# json字符串  和 class类对象参数
def praseObject(jsonStr, Class):
    """class对象需要空参构造方法,果不显示的写出构造函数，默认会自动添加一个空的构造函数 """,
    data = json.loads(jsonStr)
    result = Class()
    result.__dict__ = data
    return result


class Student(object):
    name = ''
    age = 0


student = Student()
student.name = 'jack'
student.age = 10

print(type(student.__dict__))  # {'name': 'jack', 'age': 10}
jsonStr = student.__dict__.__str__().replace("\'", "\"")
print(type(jsonStr))  # {"name": "jack", "age": 10}
newStudent = praseObject(jsonStr, Student)
print(newStudent.name)  # jack
print(newStudent.age)  # 10
