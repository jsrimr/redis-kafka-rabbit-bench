def gen():
    current = yield
    yield 10

it = gen()
print(next(it))
print(it)