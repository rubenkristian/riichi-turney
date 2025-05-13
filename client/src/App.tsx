import { lazy, Suspense, type Component } from 'solid-js';

import styles from './App.module.css';
import { Route, Router } from '@solidjs/router';

const Home = lazy(() => import('./pages/Home'));
const Player = lazy(() => import('./pages/Player'));
const Point = lazy(() => import('./pages/Point'));
const Setting = lazy(() => import('./pages/Setting'));
const Games = lazy(() => import('./pages/Games'));
const NotFound = lazy(() => import('./pages/NotFound'));

const App: Component = () => {
  return (
    <div class={styles.App}>
      <Suspense fallback={<div>Loading page...</div>}>
        <Router>
          <Route path={'/'} component={Home}/>
          <Route path={'/player'} component={Player}/>
          <Route path={'/point'} component={Point}/>
          <Route path={'/setting'} component={Setting}/>
          <Route path={'/games'} component={Games}/>
          <Route path="*404" component={NotFound} />
        </Router>
      </Suspense>
    </div>
  );
};

export default App;
