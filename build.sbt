lazy val baseSettings = Seq(
  scalaVersion := "2.13.8",
  organization := "com.github.stoneream",
  scalacOptions ++= List(
    "-Ywarn-unused",
    "-Yrangepos"
  ),
  scalafmtOnCompile := true,
  semanticdbEnabled := true,
  semanticdbVersion := scalafixSemanticdb.revision
)

lazy val root = (project in file("."))
  .settings(baseSettings: _*)
  .aggregate(diaryGenerator, migrator)

lazy val diaryGenerator = (project in file("app"))
  .settings(baseSettings: _*)
  .settings(
    libraryDependencies ++= Dependencies.bundle,
    assembly / mainClass := Some("diary_generator.Main"),
    assembly / assemblyJarName := "diary-generator.jar"
  )

lazy val migrator = (project in file("migrator"))
  .settings(baseSettings: _*)
  .settings(
    libraryDependencies ++= Dependencies.bundle,
    assembly / mainClass := Some("diary_generator.migrator.Main"),
    assembly / assemblyJarName := "diary-generator-migrator.jar"
  )
