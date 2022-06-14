from oslo_config import cfg
from oslo_log import log as logging
import sys

LOG = logging.getLogger(__name__)
CONF = cfg.CONF
DOMAIN = "logger"


def prepare_service():
    logging.register_options(CONF)
    cfg.CONF(sys.argv[1:], project=DOMAIN, default_config_files=['log.conf'])
    logging.setup(CONF, DOMAIN)


prepare_service()
