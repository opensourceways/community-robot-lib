class BaseHandler(object):
    def __init__(self, fuc=None):
        self.func = fuc

    def check_registered_func(self):
        if self.func is None:
            raise Exception('Unable to register handler because of no function. ')

    def get_func(self):
        return self.func


class PushHandler(BaseHandler):
    def __init__(self, fuc):
        super().__init__(fuc)


class TagHandler(BaseHandler):
    def __init__(self, fuc):
        super().__init__(fuc)


class IssueHandler(BaseHandler):
    def __init__(self, func):
        super().__init__(func)


class NoteHandler(BaseHandler):
    def __init__(self, fuc):
        super().__init__(fuc)


class MergeRequestHandler(BaseHandler):
    def __init__(self, fuc):
        super().__init__(fuc)


