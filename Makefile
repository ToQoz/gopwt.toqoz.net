artifact-pack: gopherjs web
	zip -r deploy.zip Dockerfile Dockerrun.aws.json web.go statics sandbox
artifact-set-to-deploy-target: artifact-pack
	ruby -r yaml -e "File.open('.elasticbeanstalk/config.yml', 'r+') { |f| data = f.read; f.seek(0, IO::SEEK_SET); f.write(YAML.load(data).tap { |y| y['deploy'] ||= {}; y['deploy']['artifact'] = 'deploy.zip' }.to_yaml) }"
deploy: artifact-set-to-deploy-target
	eb deploy
