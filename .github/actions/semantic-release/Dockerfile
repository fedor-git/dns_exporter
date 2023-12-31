FROM node:21.5.0

RUN npm install -g semantic-release@17.4.7 \
  @semantic-release/commit-analyzer@8.0.1 \
  @semantic-release/release-notes-generator@9.0.3 \
  @semantic-release/npm@7.1.3 \
  @semantic-release/changelog@5.0.1 \
  @semantic-release/git@9.0.1 \
  @semantic-release/exec \
  @saithodev/semantic-release-backmerge@1.5.3 \
  conventional-changelog-conventionalcommits@4.6.1

WORKDIR /app

ENTRYPOINT ["/usr/local/bin/npx"]
CMD ["semantic-release", "--branches", "main", "--no-ci"]
