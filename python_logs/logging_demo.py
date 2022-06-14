import logging


class LogUtils(logging.Logger):

    def __init__(self, name, log_level=logging.DEBUG,
                 log_format='[%(asctime)s][%(levelname)s][%(filename)s:%(lineno)s]%(message)s',
                 data_format='%Y-%m-%d %H:%M:%S'):
        super().__init__(name)
        formatter = logging.Formatter(log_format, data_format)
        console = logging.StreamHandler()
        console.setLevel(log_level)
        console.setFormatter(formatter)
        self.addHandler(console)
        self.setLevel(log_level)
