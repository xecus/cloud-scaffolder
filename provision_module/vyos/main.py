import sys
import click
import util
import control

logger = util.Logger.setup()


@click.group(invoke_without_command=True)
@click.pass_context
def cli(ctx):
    if ctx.invoked_subcommand is None:
        print(ctx.get_help())


@cli.command(help='provision target router')
@click.option('--target_router_ip', callback=util.validate_ip_address, required=True)
@click.option('--target_router_port', type=int, required=True)
@click.option('--target_router_username', required=True)
@click.option('--target_router_password', required=True)
@click.option('--router_hostname', required=True)
@click.option('--dhcp_ip_start', callback=util.validate_ip_address, required=True)
@click.option('--dhcp_ip_stop', callback=util.validate_ip_address, required=True)
@click.option('--dhcp_subnet', callback=util.validate_netmask, required=True)
@click.option('--dhcp_default_dns', callback=util.validate_ip_address, required=True)
@click.option('--dhcp_default_router', callback=util.validate_ip_address, required=True)
def provision(
        target_router_ip, target_router_port, target_router_username, target_router_password,
        router_hostname,
        dhcp_ip_start, dhcp_ip_stop, dhcp_subnet, dhcp_default_dns, dhcp_default_router
):

    config = {
        'target_router_ip': target_router_ip,
        'target_router_port': target_router_port,
        'target_router_username': target_router_username,
        'target_router_password': target_router_password,
        'router_hostname': router_hostname,
        'dhcp_ip_start': dhcp_ip_start,
        'dhcp_ip_stop': dhcp_ip_stop,
        'dhcp_subnet': dhcp_subnet.with_prefixlen,
        'dhcp_default_dns': dhcp_default_dns,
        'dhcp_default_router': dhcp_default_router
    }

    logger.info('boot', extra=config)

    try:
        control.provision(logger, config)
    except Exception as e:
        print(e)
        logger.error('Error: {}'.format(e))
        code = 1
    else:
        code = 0

    logger.info('finish', extra={
        'code': code
    })
    sys.exit(code)


if __name__ == "__main__":
    cli()
