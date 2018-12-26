#!/usr/bin/env ruby

#
# This script was inspired by (copied and modified from):
# http://www.darrencoxall.com/golang/cross-compiling-go-is-easy/
#

appname = (ARGV.length > 0) ? ARGV[0] : Dir.getwd.split('/').last

oses = ['linux', 'darwin', 'windows']
arches = ['amd64', '386']

Dir.mkdir('bin') if !Dir.exist?('./bin')

oses.each do |os|
  ENV['GOOS'] = os
  arches.each do |arch|
    ENV['GOARCH'] = arch
    cmd = "go build -o \"bin/#{appname}_#{ENV['GOOS']}-#{ENV['GOARCH']}\""
    system(cmd)
  end
end
