FROM node:21.5.0

RUN apt-get update && apt-get install -y git

RUN npm install -g semantic-release@22.0.12 \
  @semantic-release/commit-analyzer@11.1.0 \
  @semantic-release/release-notes-generator@12.1.0 \
  @semantic-release/npm@11.0.2 \
  @semantic-release/changelog@6.0.3 \
  @semantic-release/git@10.0.1 \
  @semantic-release/exec@6.0.3 \
  @saithodev/semantic-release-backmerge@4.0.1 \
  conventional-changelog-conventionalcommits@7.0.2

COPY package.json /github/workspace

WORKDIR /github/workspace

ENTRYPOINT ["/usr/local/bin/npx"]
CMD ["semantic-release", "--no-ci"]
