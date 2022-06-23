from robotgitee import logutil
from robotgitee import framework
from robotgitee import client

from robot import _Robot


if __name__ == "__main__":
    logutil.init_logger("test")

    cli = client.new_client("")

    bot = _Robot(cli)

    framework.run(bot, 8000, 300)
