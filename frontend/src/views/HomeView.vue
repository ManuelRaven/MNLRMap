<script setup lang="ts"></script>

<template>
  <q-page padding>
    <q-card dark bordered class="bg-grey-9 my-card">
      <q-card-section>
        <div class="text-h6">MNLR Map Tool</div>
        <div class="text-subtitle2">Map Extract</div>
      </q-card-section>

      <q-separator dark inset />

      <q-card-section>
        <h5 class="q-mt-none q-mb-md">How to Use This App</h5>

        <div class="text-body1 q-mb-md">
          This application allows you to create, manage, and use map extracts
          for offline use. Here's a comprehensive guide to all features:
        </div>

        <h6 class="q-my-md">Map Extract Tools</h6>
        <q-list bordered separator>
          <q-item>
            <q-item-section avatar>
              <q-icon name="map" color="primary" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Creating a Map Extract</q-item-label>
              <q-item-label caption>
                Navigate to the Map Extract page and use the draw rectangle tool
                located in the top-left corner of the map. Click this tool, then
                draw a rectangle on the map to define your extract area.
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item>
            <q-item-section avatar>
              <q-icon name="edit" color="primary" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Naming Your Extract</q-item-label>
              <q-item-label caption>
                After drawing a rectangle, a dialog will appear where you can
                name your extract. Enter a descriptive name and click "Create"
                to generate your map extract.
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item>
            <q-item-section avatar>
              <q-icon name="delete" color="negative" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Removing a Rectangle</q-item-label>
              <q-item-label caption>
                To remove a drawn rectangle, click on the trash icon in the
                toolbar and then select the rectangle.
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>

        <h6 class="q-my-md">Managing & Using Maps</h6>
        <q-list bordered separator>
          <q-item>
            <q-item-section avatar>
              <q-icon name="view_list" color="green" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Local Maps List</q-item-label>
              <q-item-label caption>
                Access the Local Maps List page to view all your created map
                extracts. Here you can see details about each map including
                name, creation date, and size. You can also delete maps you no
                longer need.
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item>
            <q-item-section avatar>
              <q-icon name="explore" color="green" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Local Maps Viewer</q-item-label>
              <q-item-label caption>
                Open your maps in the Local Maps Viewer to use them. This viewer
                provides tools for zooming, panning, and measuring distances.
                Your extracted maps are available offline, making them perfect
                for areas with limited connectivity.
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item>
            <q-item-section avatar>
              <q-icon name="download" color="green" />
            </q-item-section>
            <q-item-section>
              <q-item-label>PMTiles Format</q-item-label>
              <q-item-label caption>
                All map extracts are stored in the efficient PMTiles format,
                which allows for fast rendering and small file sizes. This makes
                the maps ideal for offline use on various devices.
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>

        <h6 class="q-my-md">Account & Navigation</h6>
        <q-list bordered separator>
          <q-item>
            <q-item-section avatar>
              <q-icon name="menu" color="amber" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Navigation Menu</q-item-label>
              <q-item-label caption>
                Click on the navigation drawer icon in the top-left corner of
                the app to access all available features. The menu provides
                quick access to Home, Map Extract, Local Maps List, and Local
                Maps Viewer.
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item>
            <q-item-section avatar>
              <q-icon name="person" color="amber" />
            </q-item-section>
            <q-item-section>
              <q-item-label>User Profile</q-item-label>
              <q-item-label caption>
                Access your user profile to manage account settings, update
                personal information, and customize app preferences. You can
                also view statistics about your map usage.
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>

        <div class="text-body1 q-mt-lg">
          <p>
            <strong>Note:</strong> This application is designed to work both
            online and offline. Once you've created map extracts, you can use
            them without an internet connection, making it perfect for hiking,
            travel, or any situation where connectivity might be limited.
          </p>
        </div>
      </q-card-section>
    </q-card>
  </q-page>
</template>
