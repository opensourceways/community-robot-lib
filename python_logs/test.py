import logging

from python_logs.logging_demo import LogUtils

if __name__ == '__main__':
    format_str = '[%(asctime)s]-[%(levelname)s]-[%(filename)s:%(lineno)s]%(message)s'
    logger = LogUtils('logger', logging.WARN, log_format=format_str)

    logger.info('test')
    logger.warning('test')
