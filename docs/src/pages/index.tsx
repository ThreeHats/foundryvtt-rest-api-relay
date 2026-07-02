import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '../components/HomepageFeatures';

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <h1 className="hero__title">
          {siteConfig.title}
        </h1>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/intro">
            Get Started ⚡
          </Link>
          <Link
            className="button button--outline button--secondary button--lg"
            to="/api/"
            style={{marginLeft: '10px'}}>
            API Reference 📚
          </Link>
        </div>
      </div>
    </header>
  );
}

export default function Home(): React.ReactElement {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout>
      <head>
        <title>{siteConfig.title}</title>
        <meta name="description" content="Message relay server for the foundryvtt-rest-api Foundry Module. Connect your FoundryVTT game with external applications using REST API and WebSocket connections." />
      </head>
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <section className={styles.codeSection}>
          <div className="container">
            <div className="row">
              <div className="col col--6">
                <h2>Quick Start</h2>
                <p>Get up and running with the FoundryVTT REST API Relay in minutes:</p>
                <pre className={styles.codeBlock}>
                  <code>{`# Clone and start with Docker
git clone https://github.com/JustAnotherIdea/foundryvtt-rest-api-relay.git
cd foundryvtt-rest-api-relay
docker-compose up -d

# Server available at http://localhost:3010`}</code>
                </pre>
              </div>
              <div className="col col--6">
                <h2>Example API Call</h2>
                <p>Search for entities in your FoundryVTT world:</p>
                <pre className={styles.codeBlock}>
                  <code>{`curl -X GET "http://localhost:3010/search" \\
  -H "x-api-key: your-api-key" \\
  -G -d "clientId=your-client-id" -d "query=goblin" -d "filter=Actor"`}</code>
                </pre>
              </div>
            </div>
          </div>
        </section>
      </main>
    </Layout>
  );
}
