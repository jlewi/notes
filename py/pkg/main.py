import logging
import logging.config
import os
import library

if __name__ == "__main__":
  dir_name = os.path.dirname(__file__)
  config_file = os.path.join(dir_name, 'logging.conf')
  logging.config.fileConfig(config_file)
  logger = logging.getLogger(__name__)
  logger.info('logging configured with file %s', config_file, extra={'foo': 'bar'})

  library.some_method()