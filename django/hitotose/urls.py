"""
URL configuration for hitotose project.

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/4.2/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path, include
from rest_framework import routers
from game.views import GameViewSet, badges, create_game, delete_game, get_csrf, get_games, start_game, stop_game, stopwatch, to_update_game, update_game

router = routers.DefaultRouter()
router.register(r'games', GameViewSet)

urlpatterns = [
    path('admin/', admin.site.urls),
    path('api/', include(router.urls)),
    path('api/csrf/', get_csrf, name='get_csrf'),
    path('api/game/badge', badges, name='badges'),
    path('api/game/create', create_game, name='create_game'),
    path('api/game/delete', delete_game, name='delete_game'),
    path('api/game/pages', get_games, name='get_games'),
    path('api/game/start', start_game, name='start_game'),
    path('api/game/stop/', stop_game, name='stop_game'),
    path('api/game/stopwatch', stopwatch, name='stopwatch'),
    path('api/game/update/', update_game, name='update_game'),
]
