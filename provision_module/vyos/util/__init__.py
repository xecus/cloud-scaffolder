import click
import ipaddress
import pygogo as gogo


class Logger(object):
    logger = None

    @classmethod
    def get(cls):
        if cls.logger is None:
            return cls.setup()
        return cls.logger

    @classmethod
    def setup(cls, name='default', level='info'):
        formatter = gogo.formatters.structured_formatter
        logger = gogo.Gogo(
            'struct',
            low_formatter=formatter,
            low_level=level
        ).get_logger(name)
        cls.logger = logger
        return cls.logger


def validate_ip_address(ctx, param, value):
    try:
        v = ipaddress.ip_address(value)
    except ValueError:
        raise click.BadParameter('illegal address format')
    if not isinstance(v, ipaddress.IPv4Address):
        raise click.BadParameter('only supported IPv4')
    return v


def validate_netmask(ctx, param, value):
    try:
        v = ipaddress.ip_network(value)
    except ValueError:
        raise click.BadParameter('illegal address format')
    if not isinstance(v, ipaddress.IPv4Network):
        raise click.BadParameter('only supported IPv4')
    return v
