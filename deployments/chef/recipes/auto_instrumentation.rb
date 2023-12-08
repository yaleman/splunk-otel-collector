# Cookbook:: splunk_otel_collector
# Recipe:: auto_instrumentation

with_new_instrumentation = node['splunk_otel_collector']['auto_instrumentation_version'] == 'latest' || Gem::Version.new(node['splunk_otel_collector']['auto_instrumentation_version']) >= Gem::Version.new('0.87.0')
with_systemd = node['splunk_otel_collector']['auto_instrumentation_systemd'].to_s.downcase == 'true'
with_java = node['splunk_otel_collector']['with_auto_instrumentation_sdks'].include?('java')
with_nodejs = node['splunk_otel_collector']['with_auto_instrumentation_sdks'].include?('nodejs') && with_new_instrumentation
npm_path = node['splunk_otel_collector']['auto_instrumentation_npm_path']
splunk_otel_js_path = '/usr/lib/splunk-instrumentation/splunk-otel-js.tgz'
splunk_otel_js_prefix = '/usr/lib/splunk-instrumentation/splunk-otel-js'
node.run_state[:with_nodejs] = with_nodejs && shell_out("bash -c 'command -v #{npm_path}'").exitstatus == 0

ohai 'reload packages' do
  action :nothing
  plugin 'packages'
end

directory splunk_otel_js_prefix do
  action :nothing
  recursive true
  only_if { !with_systemd && node.run_state[:with_nodejs] }
end

ruby_block 'install splunk-otel-js' do
  action :nothing
  block do
    node.run_state[:with_nodejs] = with_nodejs && shell_out("bash -c 'command -v #{npm_path}'").exitstatus == 0
    if node.run_state[:with_nodejs]
      resources(directory: splunk_otel_js_prefix).run_action(:create)
      shell_out!("#{npm_path} install --prefix #{splunk_otel_js_prefix} #{splunk_otel_js_path}")
    end
  end
  only_if { with_nodejs }
end

execute 'reload systemd' do
  action :nothing
  command 'systemctl daemon-reload'
end

package 'splunk-otel-auto-instrumentation' do
  action :install
  version node['splunk_otel_collector']['auto_instrumentation_version'] if node['splunk_otel_collector']['auto_instrumentation_version'] != 'latest'
  flush_cache [ :before ] if platform_family?('amazon', 'rhel')
  options '--allow-downgrades' if platform_family?('debian') \
    && node['packages'] \
    && node['packages']['apt'] \
    && Gem::Version.new(node['packages']['apt']['version'].split('~')[0]) >= Gem::Version.new('1.1.0')
  allow_downgrade true if platform_family?('amazon', 'rhel', 'suse')
  notifies :reload, 'ohai[reload packages]', :immediately
  notifies :run, 'ruby_block[install splunk-otel-js]', :immediately
end

file '/etc/splunk/zeroconfig/java.conf' do
  action :delete
  only_if { with_systemd || !with_java }
end

file '/etc/splunk/zeroconfig/node.conf' do
  action :delete
  only_if { with_systemd || !node.run_state[:with_nodejs] }
end

file '/usr/lib/splunk-instrumentation/instrumentation.conf' do
  action :delete
  only_if { with_systemd || with_new_instrumentation }
end

file '/usr/lib/systemd/system.conf.d/00-splunk-otel-auto-instrumentation.conf' do
  action :delete
  notifies :run, 'execute[reload systemd]', :delayed
  not_if { with_systemd }
end

template '/etc/splunk/zeroconfig/java.conf' do
  variables(
    installed_version: lazy { node['packages']['splunk-otel-auto-instrumentation']['version'] }
  )
  source 'java.conf.erb'
  only_if { with_java && with_new_instrumentation && !with_systemd }
end

template '/etc/splunk/zeroconfig/node.conf' do
  variables(
    installed_version: lazy { node['packages']['splunk-otel-auto-instrumentation']['version'] }
  )
  source 'node.conf.erb'
  only_if { node.run_state[:with_nodejs] && !with_systemd }
end

template '/usr/lib/splunk-instrumentation/instrumentation.conf' do
  variables(
    installed_version: lazy { node['packages']['splunk-otel-auto-instrumentation']['version'] }
  )
  source 'instrumentation.conf.erb'
  only_if { with_java && !with_new_instrumentation && !with_systemd }
end

directory '/usr/lib/systemd/system.conf.d' do
  recursive true
  only_if { with_systemd }
  only_if { with_java || node.run_state[:with_nodejs] }
end

template '/usr/lib/systemd/system.conf.d/00-splunk-otel-auto-instrumentation.conf' do
  variables(
    installed_version: lazy { node['packages']['splunk-otel-auto-instrumentation']['version'] },
    with_nodejs: lazy { node.run_state[:with_nodejs] }
  )
  source '00-splunk-otel-auto-instrumentation.conf.erb'
  notifies :run, 'execute[reload systemd]', :delayed
  only_if { with_systemd }
  only_if { with_java || node.run_state[:with_nodejs] }
end

template '/etc/ld.so.preload' do
  variables(
    with_systemd: lazy { with_systemd },
    with_java: lazy { with_java },
    with_nodejs: lazy { node.run_state[:with_nodejs] }
  )
  source 'ld.so.preload.erb'
end
