import vymgmt
import pexpect


class ProvisionError(Exception):
    pass


class ConfigParameterError(ProvisionError):
    pass


class ProvisionFailed(ProvisionError):
    pass


def get_vyos_connection(logger, config):
    return vymgmt.Router(
        config['target_router_ip'],
        config['target_router_username'],
        password=config['target_router_password'],
        port=config['target_router_port']
    )


def init_time_zone(vyos, logger, config):
    zone_name = 'Asia/Tokyo'
    cmd = 'system time-zone {}'.format(zone_name)
    logger.info(cmd)
    vyos.set(cmd)


def init_host_name(vyos, logger, config):
    cmd = 'system host-name {}'.format(config['router_hostname'])
    logger.info(cmd)
    vyos.set(cmd)


def init_ntp_server(vyos, logger, config):
    ntp_servers = [
        'ntp.nict.jp',
        'ntp1.jst.mfeed.ad.jp',
        'ntp2.jst.mfeed.ad.jp',
        'ntp3.jst.mfeed.ad.jp'
    ]
    for ntp_server in ntp_servers:
        cmd = 'system ntp server {}'.format(ntp_server)
        logger.info(cmd)
        vyos.set(cmd)


def init_dhcp_server(vyos, logger, config):
    dhcp_service_cmd = 'service dhcp-server shared-network-name'
    dhcp_setting_name = 'dhcp1'
    dhcp_subnet = config['dhcp_subnet']
    dhcp_ip_start = config['dhcp_ip_start']
    dhcp_ip_stop = config['dhcp_ip_stop']
    dhcp_default_router = config['dhcp_default_router']
    dhcp_default_dns = config['dhcp_default_dns']

    nat_service_cmd = 'nat source'
    nat_source_address = config['dhcp_subnet']
    nat_outbound_interface = 'eth0'
    nat_rule_id = 1

    # IP Range
    cmd = '{} {} subnet {} start {} stop {}'.format(
        dhcp_service_cmd, dhcp_setting_name, dhcp_subnet, dhcp_ip_start, dhcp_ip_stop)
    logger.info(cmd)
    vyos.set(cmd)
    # Default Router
    cmd = '{} {} subnet {} default-router {}'.format(
        dhcp_service_cmd, dhcp_setting_name, dhcp_subnet, dhcp_default_router)
    logger.info(cmd)
    vyos.set(cmd)
    # DNS Server
    cmd = '{} {} subnet {} dns-server {}'.format(
        dhcp_service_cmd, dhcp_setting_name, dhcp_subnet, dhcp_default_dns)
    logger.info(cmd)
    vyos.set(cmd)

    # Dynamic NAT
    cmd = '{} rule {} source address {}'.format(nat_service_cmd, nat_rule_id, nat_source_address)
    logger.info(cmd)
    vyos.set(cmd)
    cmd = '{} rule {} translation address masquerade'.format(nat_service_cmd, nat_rule_id)
    vyos.set(cmd)
    logger.info(cmd)
    cmd = '{} rule {} outbound-interface {}'.format(nat_service_cmd, nat_rule_id, nat_outbound_interface)
    vyos.set(cmd)
    logger.info(cmd)


def validate_provision_config(logger, config):
    required_keys = [
        'target_router_ip', 'target_router_port', 'target_router_username', 'target_router_password',
        'router_hostname',
        'dhcp_ip_start', 'dhcp_ip_stop', 'dhcp_subnet', 'dhcp_default_dns', 'dhcp_default_router'
    ]

    for required_key in required_keys:
        if required_key not in config:
            raise ConfigParameterError('key {} is not found'.format(required_key))


def provision(logger, config):

    try:
        # Validation
        validate_provision_config(logger, config)

        # Prepare
        logger.info('Establishing a connection...')
        vyos = get_vyos_connection(logger, config)
        logger.info('Logging in...')
        vyos.login()
        logger.info('Making a transition to configure mode...')
        vyos.configure()

        # System Setting
        logger.info('Setting up Host Name... ')
        init_host_name(vyos, logger, config)
        logger.info('Setting up NTP-Servers')
        init_ntp_server(vyos, logger, config)
        logger.info('Setting up TimeZone...')
        init_time_zone(vyos, logger, config)

        # DHCP Server
        logger.info('Settin up DHCP And DynamicNAT...')
        init_dhcp_server(vyos, logger, config)

        # Save
        logger.info('Commiting....')
        vyos.commit()
        logger.info('Saving....')
        vyos.save()
        logger.info('Logging out....')
        vyos.exit()
        vyos.logout()
    except ConfigParameterError as e:
        logger.error('Parameter Error: {}'.format(e))
        code = 1
    except pexpect.pxssh.ExceptionPxssh as e:
        logger.error('Shell Error: {}'.format(e))
        code = 2
    except vymgmt.router.VyOSError as e:
        logger.error('Config Error: {}'.format(e))
        code = 3
    else:
        code = 0

    if code:
        raise ProvisionFailed('provision Failed. code={}'.format(code))
